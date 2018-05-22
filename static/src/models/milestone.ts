import { Moment } from 'moment';

export interface Milestone {
  ID: number;
  IID: number;
  parallel: number;
  Title: string;
  StartDate?: Moment;
  DueDate?: Moment;
}
