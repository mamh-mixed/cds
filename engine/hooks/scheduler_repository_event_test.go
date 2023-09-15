package hooks

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/ovh/cds/engine/cache"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/cdsclient/mock_cdsclient"
	"github.com/rockbears/log"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestManageAnalysisCallback(t *testing.T) {
	log.Factory = log.NewTestingWrapper(t)
	s, cancel := setupTestHookService(t)
	defer cancel()

	event := GiteaEventPayload{}
	event.Repository.FullName = "ovh/cds"
	event.Ref = "master"
	event.After = "123456"

	bts, _ := json.Marshal(event)

	// Create event
	hr := sdk.HookRepositoryEvent{
		UUID:           sdk.UUID(),
		VCSServerName:  "private-github",
		VCSServerType:  "github",
		RepositoryName: "ovh/cds",
		Status:         sdk.HookEventStatusAnalysis,
		EventName:      "push",
		Created:        time.Now().UnixNano(),
		Body:           bts,
		Analyses: []sdk.HookRepositoryEventAnalysis{
			{
				ProjectKey: "MYPROJECT",
				Status:     sdk.RepositoryAnalysisStatusInProgress,
				AnalyzeID:  sdk.UUID(),
			},
		},
	}
	require.NoError(t, s.Dao.SaveRepositoryEvent(context.TODO(), &hr))

	// Create repo
	repoKey := s.Dao.GetRepositoryMemberKey(hr.VCSServerType, hr.VCSServerName, hr.RepositoryName)
	_, err := s.Dao.CreateRepository(context.TODO(), repoKey, hr.VCSServerType, hr.VCSServerName, hr.RepositoryName)
	require.NoError(t, err)

	callback := sdk.HookAnalysisCallback{
		AnalysisID:     hr.Analyses[0].AnalyzeID,
		RepositoryName: hr.RepositoryName,
		VCSServerName:  hr.VCSServerName,
		VCSServerType:  hr.VCSServerType,
		HookEventUUID:  hr.UUID,
		AnalysisStatus: sdk.RepositoryAnalysisStatusSucceed,
	}

	require.NoError(t, s.updateHookEventWithCallback(context.TODO(), callback))

	k := cache.Key(repositoryEventRootKey, s.Dao.GetRepositoryMemberKey(hr.VCSServerType, hr.VCSServerName, hr.RepositoryName), hr.UUID)
	var hreUpdate sdk.HookRepositoryEvent
	f, err := s.Cache.Get(k, &hreUpdate)
	require.NoError(t, err)
	require.True(t, f)
	require.Equal(t, sdk.RepositoryAnalysisStatusSucceed, hreUpdate.Analyses[0].Status)

}

func TestManageRepositoryEvent_PushEventTriggerAnalysis(t *testing.T) {
	log.Factory = log.NewTestingWrapper(t)
	s, cancel := setupTestHookService(t)
	defer cancel()

	event := GiteaEventPayload{}
	event.Repository.FullName = "ovh/cds"
	event.Ref = "master"
	event.After = "123456"

	bts, _ := json.Marshal(event)

	// Create event
	hr := sdk.HookRepositoryEvent{
		UUID:           sdk.UUID(),
		VCSServerName:  "private-github",
		VCSServerType:  "github",
		RepositoryName: "ovh/cds",
		Status:         sdk.HookEventStatusScheduled,
		EventName:      "push",
		Created:        time.Now().UnixNano(),
		Body:           bts,
	}
	require.NoError(t, s.Dao.SaveRepositoryEvent(context.TODO(), &hr))

	// Create repo
	repoKey := s.Dao.GetRepositoryMemberKey(hr.VCSServerType, hr.VCSServerName, hr.RepositoryName)
	_, err := s.Dao.CreateRepository(context.TODO(), repoKey, hr.VCSServerType, hr.VCSServerName, hr.RepositoryName)
	require.NoError(t, err)

	s.Client.(*mock_cdsclient.MockInterface).EXPECT().HookRepositoriesList(gomock.Any(), hr.VCSServerName, hr.RepositoryName).Return([]sdk.ProjectRepository{
		{
			Name:       hr.RepositoryName,
			ProjectKey: "TEST",
		},
	}, nil).Times(1)
	s.Client.(*mock_cdsclient.MockInterface).EXPECT().ProjectRepositoryAnalysis(gomock.Any(), gomock.Any()).Times(1)

	// Force dequeue
	k := cache.Key(repositoryEventRootKey, s.Dao.GetRepositoryMemberKey(hr.VCSServerType, hr.VCSServerName, hr.RepositoryName), hr.UUID)
	require.NoError(t, s.manageRepositoryEvent(context.TODO(), k))
}

func TestManageRepositoryEvent_NonPushEventNoTriggerAnalysis(t *testing.T) {
	log.Factory = log.NewTestingWrapper(t)
	s, cancel := setupTestHookService(t)
	defer cancel()

	event := GiteaEventPayload{}
	event.Repository.FullName = "ovh/cds"
	event.Ref = "master"
	event.After = "123456"

	bts, _ := json.Marshal(event)

	// Create event
	hr := sdk.HookRepositoryEvent{
		UUID:           sdk.UUID(),
		VCSServerName:  "private-github",
		VCSServerType:  "github",
		RepositoryName: "ovh/cds",
		Status:         sdk.HookEventStatusScheduled,
		EventName:      "pull_request",
		Created:        time.Now().UnixNano(),
		Body:           bts,
	}
	require.NoError(t, s.Dao.SaveRepositoryEvent(context.TODO(), &hr))

	// Create repo
	repoKey := s.Dao.GetRepositoryMemberKey(hr.VCSServerType, hr.VCSServerName, hr.RepositoryName)
	_, err := s.Dao.CreateRepository(context.TODO(), repoKey, hr.VCSServerType, hr.VCSServerName, hr.RepositoryName)
	require.NoError(t, err)

	// Force dequeue
	k := cache.Key(repositoryEventRootKey, s.Dao.GetRepositoryMemberKey(hr.VCSServerType, hr.VCSServerName, hr.RepositoryName), hr.UUID)
	require.NoError(t, s.manageRepositoryEvent(context.TODO(), k))
}
