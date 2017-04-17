import axios from 'axios';
import { apiAddress } from './config';

export default {
  getTitle(titleId) {
    return axios.get(`${apiAddress}/api/titles/${titleId}`);
  },

  listTitles() {
    return axios.get(`${apiAddress}/api/titles`);
  },

  importTitles(file) {
    const formData = new FormData();
    formData.append('file', file, file.name);

    return axios.post(`${apiAddress}/api/import`, formData);
  },
};
