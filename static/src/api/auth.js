import axios from 'axios';
import { apiAddress } from './config';

export default {
  requestAuth() {
    return axios.get(`${apiAddress}/api/auth`);
  },

  processAuth(state, code) {
    return axios.post(`${apiAddress}/api/auth/callback`, { state, code });
  },
};
