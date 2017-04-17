import axios from 'axios';
import { apiAddress } from './config';

export default {
  getImports(titleId, rplId) {
    return axios.get(
      `${apiAddress}/api/titles/${titleId}/rpls/${rplId}/imports`,
    );
  },
};
