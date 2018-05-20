import axios from 'axios';
import { Work } from '../models/work';
import moment = require('moment');
import { Milestone } from '../models/milestone';

export class Api {
  static fetchWorks(): Promise<Work[]> {
    return axios.post('/api/works').then(r => {
      const works = r.data;
      return works.map((w: any) => {
        w.Issue.DueDate = w.Issue.DueDate ? moment(w.Issue.DueDate) : null;
        return w;
      });
    });
  }

  static fetchMilestones(): Promise<Milestone[]> {
    return axios.post('/api/milestones').then(r => {
      const milestones = r.data;
      return milestones.map((m: any) => {
        m.StartDate = m.StartDate ? moment(m.StartDate) : null;
        m.DueDate = m.DueDate ? moment(m.DueDate) : null;
        return m;
      });
    });
  }
}
