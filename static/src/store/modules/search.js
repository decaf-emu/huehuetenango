import * as types from '../mutation-types';
import search from '../../api/search';

const state = {
  titleResults: [],
};

const getters = {
  titleSearchResults: state => state.titleResults,
};

const actions = {
  searchTitles({ commit }, term) {
    search.searchTitles(
      term,
      results => commit(types.SEARCH_TITLES_SUCCESS, { results }),
      () => commit(types.SEARCH_TITLES_FAILURE),
    );
  },
};

const mutations = {
  [types.SEARCH_TITLES_SUCCESS](state, { results }) {
    state.titleResults = results;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};
