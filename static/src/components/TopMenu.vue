<template>
  <nav class="uk-navbar-container uk-margin" uk-navbar>
    <div class="nav-content uk-navbar-left">
      <a class="uk-navbar-toggle" uk-navbar-toggle-icon uk-toggle="target: #offcanvas-nav-primary" href="#"></a>
      <router-link class="uk-navbar-item uk-logo" :to="{ name: 'titles' }">Huehuetenango</router-link>
    </div>
    <div class="uk-navbar-left uk-flex-1">
      <div class="uk-navbar-item uk-width-expand">
        <form class="uk-search uk-search-navbar uk-width-1-1">
          <span uk-search-icon></span>
          <input class="uk-search-input" type="search" placeholder="Search..." v-model="searchTerm">
          <div v-if="searchTerm && searchTerm.length > 0">
            <vk-dropdown position="bottom center" class="uk-width-xxxlarge uk-margin uk-margin-top" :scrollable="true"
              offset="-50px -50px" :show="true">
              <SearchResults />
            </vk-dropdown>
          </div>
        </form>
      </div>
    </div>
    <div class="uk-navbar-right uk-margin uk-margin-right">
      <router-link :to="{ name: 'login' }" class="uk-button uk-button-primary" v-if="!isLoggedIn">
        Sign in with GitHub
      </router-link>
    </div>
    <ul class="uk-navbar-nav uk-margin-small-right" v-if="isLoggedIn">
      <li>
        <a href="#">
          <span class="uk-icon uk-icon-image uk-margin-small-right" :style="{ backgroundImage: 'url(' + avatarUrl + ')' }"></span>
          {{ name }}
        </a>
        <div class="uk-navbar-dropdown">
          <ul class="uk-nav uk-navbar-dropdown-nav">
            <li><router-link :to="{ name: 'logout' }">Logout</router-link></li>
          </ul>
        </div>
      </li>
    </ul>
  </nav>
</template>

<script>
import { mapGetters } from 'vuex';
import SearchResults from './search/SearchResults.vue';

export default {
  components: { SearchResults },
  data() {
    return {
      searchTerm: null,
    };
  },
  computed: {
    ...mapGetters(['isLoggedIn', 'name', 'avatarUrl']),
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
