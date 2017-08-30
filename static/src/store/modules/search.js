import * as types from '../mutation-types';
import search from '../../api/search';

export default {
  state: {
    term: null,
    titleResults: [],
    functionResults: [],
    dataResults: [],
  },

  getters: {
    searchTerm: state => state.term,
    titleSearchResults: state => state.titleResults,
    functionSearchResults: state => state.functionResults,
    dataSearchResults: state => state.dataResults,
  },

  actions: {
    async search({ commit }, term) {
      commit(types.SEARCH_TERM, { term });

      try {
        const response = await search.perform(term);
        commit(types.SEARCH_SUCCESS, {
          titles: response.data && response.data.titles ? response.data.titles : [],
          exports: response.data && response.data.exports ? response.data.exports : [],
        });
      } catch (error) {
        commit(types.SEARCH_FAILURE, { error });
      }
    },

    clearSearch({ commit }) {
      commit(types.SEARCH_TERM, { term: null });
    },
  },

  mutations: {
    [types.SEARCH_TERM](state, { term }) {
      state.term = term;
    },
    [types.SEARCH_SUCCESS](state, { titles, exports }) {
      state.titleResults = titles;
      state.functionResults = exports.filter(item => item.Type === 'func');
      state.dataResults = exports.filter(item => item.Type === 'data');
    },
  },
};
