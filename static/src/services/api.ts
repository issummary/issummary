import axios from 'axios';
import moment = require('moment');
import { IMilestone } from '../models/milestone';
import { IWork } from '../models/work';

export class Api {
  public static fetchWorks(): Promise<IWork[]> {
    return axios.post('/api/works').then(r => {
      if (r.data.hasOwnProperty('Error')) {
        return Promise.reject(r.data);
      }

      const works = r.data;
      return works.map((w: any) => {
        w.Issue.DueDate = w.Issue.DueDate ? moment(w.Issue.DueDate) : null;
        return w;
      });
    });
  }

  public static fetchMilestones(): Promise<IMilestone[]> {
    return axios.post('/api/milestones').then(r => {
      const milestones = r.data;

      if (!milestones) {
        return [];
      }

      return milestones.map((m: any) => {
        m.StartDate = m.StartDate ? moment(m.StartDate) : null;
        m.DueDate = m.DueDate ? moment(m.DueDate) : null;
        return m;
      });
    });
  }
}
