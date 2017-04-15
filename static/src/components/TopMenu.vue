<template>
  <nav class="uk-navbar-container uk-margin" uk-navbar>
    <div class="nav-content uk-navbar-left">
      <router-link class="uk-navbar-item uk-logo" :to="{ name: 'titles' }">Huehuetenango</router-link>
      <ul class="uk-navbar-nav">
        <li>
          <router-link :to="{ name: 'titles' }">Titles</router-link>
        </li>
        <li>
          <router-link :to="{ name: 'title',  params: { titleId: '000500101000400A' }}">System Library</router-link>
        </li>
      </ul>
    </div>

    <div class="uk-navbar-right">
      <form class="uk-search uk-search-navbar">
        <span uk-search-icon></span>
        <input class="uk-search-input" type="search" placeholder="Search..." v-model="searchTerm">
        <div v-if="searchTerm && searchTerm.length > 0">
          <vk-dropdown position="top right" class="uk-width-xlarge" :scrollable="true" offset="100px 0px" :show="true">
            <SearchResults />
          </vk-dropdown>
        </div>
      </form>
    </div>
  </nav>
</template>

<script>
import SearchResults from './search/SearchResults.vue';

export default {
  components: { SearchResults },
  data() {
    return {
      searchTerm: null,
    };
  },
  watch: {
    searchTerm(term) {
      this.$store.dispatch('searchTitles', term);
    },
  },
};
</script>

<style>
.navbar-search-form {
  margin-bottom: 0;
}
</style>
