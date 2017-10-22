import "rxjs/add/operator/toPromise";
import { Injectable } from "@angular/core";
import { AgentStatus, Filter } from "../model/agent";
import { Http } from "@angular/http";

@Injectable()
export class AgentService {
  private agenturl = "http://45.76.163.62:9000/v1/agents"; // URL to web api
  private filterurl = "http://45.76.163.62:/v1/testfilter";
  constructor(private http: Http) {}

  getAgents(): Promise<AgentStatus[]> {
    return this.http
      .get(this.agenturl)
      .toPromise()
      .then(response => response.json().status as AgentStatus[])
      .catch(this.handleError);
  }

  getFilter(): Promise<Filter> {
    return this.http
      .get(this.filterurl)
      .toPromise()
      .then(response => response.json() as Filter)
      .catch(this.handleError);
  }

  postFilter(filter: Filter): Promise<void> {
    return this.http
      .post(this.filterurl, filter)
      .toPromise()
      .then(response => response.json() as Object)
      .catch(this.handleError);
  }

  private handleError(error: any): Promise<any> {
    console.error("An error occurred", error); // for demo purposes only
    return Promise.reject(error.message || error);
  }
}
