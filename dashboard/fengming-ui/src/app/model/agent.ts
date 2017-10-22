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

export class Filter {
  Repo: string;
  AgentName: string;
}
