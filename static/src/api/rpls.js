import axios from 'axios';
import { apiAddress } from './config';

export default {
  async getTitleRpls(titleId) {
    return await axios.get(`${apiAddress}/api/titles/${titleId}/rpls`);
  },

  async getRpl(titleHexID, id) {
    return await axios.get(`${apiAddress}/api/titles/${titleHexID}/rpls/${id}`);
  },
};
