import axios from 'axios';
import { apiAddress } from './config';

export default {
  searchTitles(term) {
    return axios.post(`${apiAddress}/api/search`, { term });
  },
};
