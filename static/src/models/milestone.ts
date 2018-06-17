import { Moment } from 'moment';

export interface IMilestone {
  ID: number;
  IID: number;
  parallel: number;
  Title: string;
  StartDate?: Moment;
  DueDate?: Moment;
}
