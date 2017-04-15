import Vue from 'vue';
import { apiAddress } from './config';

export default {
  getTitle(hexId, cb, errCb) {
    return Vue.http.get(`${apiAddress}/api/titles/${hexId}`).then(
      response => {
        cb(response.body);
      },
      response => {
        errCb(response.body);
      },
    );
  },

  listTitles(cb, errCb) {
    return Vue.http.get(`${apiAddress}/api/titles`).then(
      response => {
        cb(response.body);
      },
      response => {
        errCb(response.body);
      },
    );
  },

  importTitles(file, cb, errCb) {
    const formData = new FormData();
    formData.append('file', file, file.name);

    return Vue.http.post(`${apiAddress}/api/import`, formData).then(
      response => {
        cb(response.body);
      },
      response => {
        errCb(response.body);
      },
    );
  },
};
