import axios from 'axios';

export class Api {
  static fetchWorks(): any {
    return axios.post('/works').then(r => r.data);
  }
}
