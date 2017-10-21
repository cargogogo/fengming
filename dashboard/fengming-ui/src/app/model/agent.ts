// AgentStatus ...
export class AgentStatus {
  Name: string;
  Addr: string;
  Tasks: Task[];
}

export class Task {
  ID: string;
  LayerName: string;
  Status: string;
}
