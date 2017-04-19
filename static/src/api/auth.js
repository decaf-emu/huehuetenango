import axios from 'axios';
import { apiAddress } from './config';

export default {
  async requestAuth() {
    return await axios.get(`${apiAddress}/api/auth`);
  },

  async processAuth(state, code) {
    return await axios.post(`${apiAddress}/api/auth/callback`, { state, code });
  },
};
