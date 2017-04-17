import * as types from '../mutation-types';
import titles from '../../api/titles';

const state = {
  title: null,
  all: [],
  importing: false,
  importSuccess: false,
  importFailed: false,
};

const getters = {
  allTitles: state => state.all,
  title: state => state.title,
  importingTitles: state => state.importing,
  importTitlesSuccess: state => state.importSuccess,
  importTitlesFailed: state => state.importFailed,
};

const actions = {
  getAllTitles({ commit }) {
    titles
      .listTitles()
      .then(({ data }) => commit(types.ALL_TITLES_SUCCESS, { titles: data }))
      .catch(() => commit(types.ALL_TITLES_FAILURE));
  },

  getTitle({ commit }, titleId) {
    commit(types.GET_TITLE_LOADING);

    titles
      .getTitle(titleId)
      .then(({ data }) => commit(types.GET_TITLE_SUCCESS, { title: data }))
      .catch(() => commit(types.GET_TITLE_FAILURE));
  },

  importTitles({ commit }, file) {
    commit(types.IMPORT_TITLES_LOADING);

    titles.importTitles(
      file,
      () => commit(types.IMPORT_TITLES_SUCCESS),
      () => commit(types.IMPORT_TITLES_FAILURE),
    );
  },
};

const mutations = {
  [types.ALL_TITLES_SUCCESS](state, { titles }) {
    state.all = titles.sort((a, b) => {
      const nameA = a.LongNameEnglish.toUpperCase();
      const nameB = b.LongNameEnglish.toUpperCase();

      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }

      return 0;
    });
  },

  [types.GET_TITLE_LOADING](state) {
    state.title = null;
  },
  [types.GET_TITLE_SUCCESS](state, { title }) {
    state.title = title;
  },

  [types.IMPORT_TITLES_LOADING](state) {
    state.importing = true;
    state.importSuccess = false;
    state.importFailed = false;
  },
  [types.IMPORT_TITLES_SUCCESS](state) {
    state.importing = false;
    state.importSuccess = true;
    state.importFailed = false;
  },
  [types.IMPORT_TITLES_FAILURE](state) {
    state.importing = false;
    state.importSuccess = false;
    state.importFailed = true;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};
