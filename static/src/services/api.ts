import axios from 'axios';
import { Work } from '../models/work';
import moment = require('moment');

export class Api {
  static fetchWorks(): Promise<Work[]> {
    return axios.post('/works').then(r => {
      const works = r.data;
      return works.map((w: any) => {
        w.Issue.DueDate = w.Issue.DueDate ? moment(w.Issue.DueDate) : null;
        return w;
      });
    });
  }
}
