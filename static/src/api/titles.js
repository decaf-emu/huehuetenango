import axios from 'axios';
import { apiAddress } from './config';

export default {
  async getTitle(titleId) {
    return await axios.get(`${apiAddress}/api/titles/${titleId}`);
  },

  async listTitles() {
    return await axios.get(`${apiAddress}/api/titles`);
  },

  async importTitles(file) {
    const formData = new FormData();
    formData.append('file', file, file.name);

    return await axios.post(`${apiAddress}/api/import`, formData);
  },
};
