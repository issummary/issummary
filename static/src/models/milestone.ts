import { Moment } from 'moment';

export interface Milestone {
  ID: number;
  IID: number;
  Title: string;
  StartDate?: Moment;
  DueDate?: Moment;
}
