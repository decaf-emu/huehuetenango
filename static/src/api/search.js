import Vue from 'vue';
import { apiAddress } from './config';

export default {
  searchTitles(term, cb, errCb) {
    const formData = new FormData();
    formData.append('term', term);

    return Vue.http.post(`${apiAddress}/api/search`, formData).then(
      response => {
        cb(response.body);
      },
      response => {
        errCb(response.body);
      },
    );
  },
};
