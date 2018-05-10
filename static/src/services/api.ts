import axios from 'axios';
import { Work } from '../models/work';

export class Api {
  static fetchWorks(): Promise<Work[]> {
    return axios.post('/works').then(r => r.data);
  }
}
