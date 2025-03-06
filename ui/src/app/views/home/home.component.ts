import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from "@angular/core";
import { Bookmark, BookmarkType } from "app/model/bookmark.model";
import { UserService } from "app/service/services.module";
import { AutoUnsubscribe } from "app/shared/decorator/autoUnsubscribe";
import { lastValueFrom, Subscription } from "rxjs";

@Component({
	selector: 'app-home',
	templateUrl: './home.html',
	styleUrls: ['./home.scss'],
	changeDetection: ChangeDetectionStrategy.OnPush
})
@AutoUnsubscribe()
export class HomeComponent implements OnInit, OnDestroy {

	bookmarks: Array<Bookmark> = [];
	projectsSubscription: Subscription;
	workflowsSubscription: Subscription;
	recentItems: Array<any> = [];
	loading: boolean;

	constructor(
		private _cd: ChangeDetectorRef,
		private _userService: UserService
	) { }

	ngOnDestroy(): void { } // Should be set to use @AutoUnsubscribe with AOT

	ngOnInit(): void {
		this.load();
	}

	async load() {
		this.loading = true;
		this._cd.markForCheck();
		this.bookmarks = await lastValueFrom(this._userService.getBookmarks());
		this.loading = false;
		this._cd.markForCheck();
	}

	generateBookmarkLink(b: Bookmark): Array<string> {
		const splitted = b.id.split('/');
		switch (b.type) {
			case BookmarkType.Workflow:
				const project = splitted.shift();
				return ['/project', project, 'run'];
			case BookmarkType.WorkflowLegacy:
				return ['/project', splitted[0], 'workflow', splitted[1]];
			case BookmarkType.Project:
				return ['/project', b.id];
			default:
				return [];
		}
	}

	generateBookmarkQueryParams(b: Bookmark, variant?: string): any {
		const splitted = b.id.split('/');
		switch (b.type) {
			case BookmarkType.Workflow:
				splitted.shift();
				const workflow_path = splitted.join('/');
				let params = { workflow: workflow_path };
				if (variant) {
					params['ref'] = variant;
				}
				return params;
			default:
				return {};
		}
	}

}