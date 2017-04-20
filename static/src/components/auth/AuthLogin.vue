<template>
  <div class="uk-overlay uk-overlay-primary uk-position-cover">
    <div class="uk-position-center" uk-spinner></div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import URI from 'urijs';

export default {
  computed: mapGetters({
    redirectUrl: 'authRedirectUrl',
  }),

  methods: {
    ...mapActions(['requestAuth']),
  },

  watch: {
    redirectUrl(url) {
      const callbackRoute = this.$router.resolve({ name: 'login-callback' });
      let callbackUri = URI(callbackRoute.href).absoluteTo(window.location);

      const authUri = URI(url).addQuery('redirect_uri', callbackUri);
      window.location = authUri.href();
    },
  },

  beforeMount() {
    this.requestAuth();
  },
};
</script>

<style>
</style>
