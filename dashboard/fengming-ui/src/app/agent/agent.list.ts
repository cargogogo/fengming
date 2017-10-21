import { Component, OnInit, ViewEncapsulation, ViewChild } from "@angular/core";
import { AgentStatus } from "../model/agent";
import { AgentService } from "../service/agent";
import { Router } from "@angular/router";

@Component({
  selector: "app-agent-list",
  templateUrl: "./agent.list.html"
})
export class AgentListComponent implements OnInit {
  @ViewChild("myTable") table: any;
  agents: AgentStatus[] = [];
  expand = false;
  constructor(private agentService: AgentService, private router: Router) {}
  ngOnInit(): void {
    this.agentService.getAgents().then(agents => {
      this.agents = agents;
      console.log(this.agents);
    });
  }

  getRowHeight(row): number {
    var h: number = row["Tasks"].length * 100;
    // console.log('getDetailHeight', h)
    return h;
  }

  toggleExpandRow(row) {
    console.log("Toggled Expand Row!", row, this.expand);
    this.expand = !this.expand;
    this.table.rowDetail.toggleExpandRow(row);
  }

  onDetailToggle(event) {
    console.log("Detail Toggled", event);
  }
}
