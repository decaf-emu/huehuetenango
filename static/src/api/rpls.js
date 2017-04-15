import Vue from 'vue';
import { apiAddress } from './config';

export default {
  getTitleRpls(hexID, cb, errCb) {
    return Vue.http.get(`${apiAddress}/api/titles/${hexID}/rpls`).then(
      response => {
        cb(response.body);
      },
      response => {
        errCb(response.body);
      },
    );
  },
  getRpl(titleHexID, id, cb, errCb) {
    return Vue.http
      .get(`${apiAddress}/api/titles/${titleHexID}/rpls/${id}`)
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
