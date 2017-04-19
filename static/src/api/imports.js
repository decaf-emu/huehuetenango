import axios from 'axios';
import { apiAddress } from './config';

export default {
  async getImports(titleId, rplId) {
    return await axios.get(
      `${apiAddress}/api/titles/${titleId}/rpls/${rplId}/imports`,
    );
  },
};
