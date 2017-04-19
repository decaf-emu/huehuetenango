import * as types from '../mutation-types';
import titles from '../../api/titles';

export default {
  state: {
    all: [],
    loadingAll: false,

    title: null,
    loadingTitle: false,

    importing: false,
    importSuccess: false,
    importFailed: false,
  },

  getters: {
    allTitles: state => state.all,
    loadingAllTitles: state => state.loadingAll,

    title: state => state.title,
    loadingTitle: state => state.loadingTitle,

    importingTitles: state => state.importing,
    importTitlesSuccess: state => state.importSuccess,
    importTitlesFailed: state => state.importFailed,
  },

  actions: {
    async getAllTitles({ commit }) {
      commit(types.ALL_TITLES_LOADING);

      try {
        const response = await titles.listTitles();
        commit(types.ALL_TITLES_SUCCESS, {
          titles: response.data,
        });
      } catch (error) {
        commit(types.ALL_TITLES_FAILURE, { error });
      }
    },

    async getTitle({ commit }, titleId) {
      commit(types.GET_TITLE_LOADING);

      try {
        const response = await titles.getTitle(titleId);
        commit(types.GET_TITLE_SUCCESS, {
          title: response.data,
        });
      } catch (error) {
        commit(types.GET_TITLE_FAILURE, { error });
      }
    },

    clearTitle({ commit }) {
      commit(types.CLEAR_TITLE);
    },

    async importTitles({ commit }, file) {
      commit(types.IMPORT_TITLES_LOADING);

      try {
        await titles.importTitles(file);
        commit(types.IMPORT_TITLES_SUCCESS);
      } catch (error) {
        commit(types.IMPORT_TITLES_FAILURE, { error });
      }
    },
  },

  mutations: {
    [types.ALL_TITLES_LOADING](state) {
      state.loadingAll = true;
    },

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

      state.loadingAll = false;
    },

    [types.ALL_TITLES_FAILED](state) {
      state.loadingAll = false;
    },

    [types.GET_TITLE_LOADING](state) {
      state.loadingTitle = true;
    },

    [types.GET_TITLE_SUCCESS](state, { title }) {
      state.title = title;
      state.loadingTitle = false;
    },

    [types.GET_TITLE_FAILED](state) {
      state.loadingTitle = false;
    },

    [types.CLEAR_TITLE](state) {
      state.title = null;
      state.loadingTitle = false;
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
  },
};
