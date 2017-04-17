import axios from 'axios';
import { apiAddress } from './config';

export default {
  getExports(titleId, rplId) {
    return axios.get(
      `${apiAddress}/api/titles/${titleId}/rpls/${rplId}/exports`,
    );
  },
};
