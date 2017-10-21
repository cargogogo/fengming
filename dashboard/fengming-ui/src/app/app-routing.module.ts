import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { AgentListComponent } from "./agent/agent.list";

const routes: Routes = [
  { path: "", redirectTo: "/dashboard", pathMatch: "full" },
  { path: "dashboard", component: AgentListComponent } // agentlist
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
