import Vue from 'vue';
import { apiAddress } from './config';

export default {
  getExports(titleHexID, rplId, cb, errCb) {
    return Vue.http
      .get(`${apiAddress}/api/titles/${titleHexID}/rpls/${rplId}/exports`)
      .then(
        response => {
          cb(response.body);
        },
        response => {
          errCb(response.body);
        },
      );
  },
};
