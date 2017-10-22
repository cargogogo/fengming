import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { NgbModule } from "@ng-bootstrap/ng-bootstrap";
import { AppComponent } from "./app.component";
import { HttpModule } from "@angular/http";
import { FormsModule } from "@angular/forms";
import { AgentService } from "./service/agent";
import { AgentListComponent } from "./agent/agent.list";
import { NgxDatatableModule } from "@swimlane/ngx-datatable";
import { AppRoutingModule } from "./app-routing.module";
import { MdDialogModule } from "@angular/material";
import {
  MdButtonModule,
  MdToolbarModule,
  MdSidenavModule,
  MdListModule
} from "@angular/material";
import { MdInputModule, MdCheckboxModule } from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { MdCardModule } from "@angular/material";

// In your App's module:

@NgModule({
  declarations: [AppComponent, AgentListComponent],
  imports: [
    NgbModule.forRoot(),
    BrowserModule,
    AppRoutingModule,
    HttpModule,
    FormsModule,
    MdButtonModule,
    MdToolbarModule,
    MdSidenavModule,
    MdDialogModule,
    MdInputModule,
    MdCheckboxModule,
    MdListModule,
    MdCardModule,
    NgxDatatableModule,
    BrowserAnimationsModule
  ],
  providers: [AgentService],
  bootstrap: [AppComponent]
})
export class AppModule {}
