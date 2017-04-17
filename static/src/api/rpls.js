import axios from 'axios';
import { apiAddress } from './config';

export default {
  getTitleRpls(titleId) {
    return axios.get(`${apiAddress}/api/titles/${titleId}/rpls`);
  },
  getRpl(titleHexID, id) {
    return axios.get(`${apiAddress}/api/titles/${titleHexID}/rpls/${id}`);
  },
};
