import * as types from '../mutation-types';
import search from '../../api/search';

export default {
  state: {
    titleResults: [],
  },

  getters: {
    titleSearchResults: state => state.titleResults,
  },

  actions: {
    async searchTitles({ commit }, term) {
      try {
        const response = await search.searchTitles(term);
        commit(types.SEARCH_TITLES_SUCCESS, {
          results: response.data,
        });
      } catch (error) {
        commit(types.SEARCH_TITLES_FAILURE, { error });
      }
    },
  },

  mutations: {
    [types.SEARCH_TITLES_SUCCESS](state, { results }) {
      state.titleResults = results;
    },
  },
};
