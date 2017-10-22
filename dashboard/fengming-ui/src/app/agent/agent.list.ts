import {
  Component,
  OnInit,
  Inject,
  ViewEncapsulation,
  ViewChild
} from "@angular/core";
import { AgentStatus, Filter } from "../model/agent";
import { AgentService } from "../service/agent";
import { Router } from "@angular/router";
import { MdDialog, MdDialogRef, MD_DIALOG_DATA } from "@angular/material";

@Component({
  selector: "app-agent-list",
  templateUrl: "./agent.list.html"
})
export class AgentListComponent implements OnInit {
  @ViewChild("myTable") table: any;
  agents: AgentStatus[] = [];
  expand = false;
  filter: Filter;
  refreshOn: Boolean;
  timer;
  constructor(
    private agentService: AgentService,
    private router: Router,
    public dialog: MdDialog
  ) {}
  ngOnInit(): void {
    this.refresh();
    this.agentService.getAgents().then(agents => {
      this.agents = agents;
      console.log(this.agents);
    });
    this.agentService.getFilter().then(filter => {
      this.filter = filter;
      console.log(this.filter);
    });
  }

  ngOnDestroy(): void {
    console.log("OnDestroy");
    this.stopRefresh();
  }
  stopRefresh(): void {
    this.refreshOn = false;
    clearInterval(this.timer);
  }

  refresh(): void {
    this.refreshOn = true;
    this.timer = setInterval(() => {
      this.agentService.getAgents().then(agents => {
        this.agents = agents;
        console.log(this.agents);
      });
    }, 3000);
  }

  getRowHeight(row): number {
    var h: number = row["Tasks"].length * 30 + 200;
    // console.log('getDetailHeight', h)
    return h;
  }

  toggleExpandRow(row) {
    console.log("Toggled Expand Row!", row, row.expand);
    row.expand = !row.expand;
    this.table.rowDetail.toggleExpandRow(row);
  }

  onDetailToggle(event) {
    console.log("Detail Toggled", event);
  }

  openDialog(): void {
    let dialogRef = this.dialog.open(DialogOverviewExampleDialog, {
      width: "250px",
      height: "250px",
      data: { fiter: this.filter }
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log("The dialog was closed");
      this.filter = result;
    });
  }
}

@Component({
  selector: "dialog-filter-dialog",
  templateUrl: "dialog-filter.html",
  styleUrls: ["./agent.css"]
})
export class DialogOverviewExampleDialog {
  constructor(
    public dialogRef: MdDialogRef<DialogOverviewExampleDialog>,
    @Inject(MD_DIALOG_DATA) public data: any
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}
