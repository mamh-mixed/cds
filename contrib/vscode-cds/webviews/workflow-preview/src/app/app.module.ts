import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzAlertModule } from 'ng-zorro-antd/alert';
import { NzFormModule } from 'ng-zorro-antd/form';
import { WorkflowGraphModule} from 'workflow-graph';
import { AppComponent } from './app.component';
import { FormsModule } from '@angular/forms';
import { CaretRightOutline, CaretDownOutline } from '@ant-design/icons-angular/icons';
import { IconDefinition } from '@ant-design/icons-angular';

const icons: IconDefinition[] = [ CaretRightOutline, CaretDownOutline ];
@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    NzIconModule.forRoot(icons),
    NzAlertModule,
    NzFormModule,
    WorkflowGraphModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
