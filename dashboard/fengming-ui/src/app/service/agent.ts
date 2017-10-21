import "rxjs/add/operator/toPromise";
import { Injectable } from "@angular/core";
import { AgentStatus } from "../model/agent";
import { Http } from "@angular/http";

@Injectable()
export class AgentService {
  private agenturl = "http://127.0.0.1:7100/v1/test"; // URL to web api
  constructor(private http: Http) {}

  getAgents(): Promise<AgentStatus[]> {
    return this.http
      .get(this.agenturl)
      .toPromise()
      .then(response => response.json() as AgentStatus[])
      .catch(this.handleError);
  }

  private handleError(error: any): Promise<any> {
    console.error("An error occurred", error); // for demo purposes only
    return Promise.reject(error.message || error);
  }
}
