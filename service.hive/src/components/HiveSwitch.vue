<template>
  <div>
    <div class="icons">
      <span v-if="error"><font-awesome-icon :icon="['fas', 'exclamation-triangle']" /></span>
      <span v-else><font-awesome-icon :icon="['far', 'check-circle']" /></span>
      <font-awesome-icon :icon="['fas', 'random']" size="lg" />
      <button v-if="asset.state === 'on'" v-on:click="toggle"><font-awesome-icon :icon="['fas', 'toggle-on']" size="lg" /></button>
      <button v-else v-on:click="toggle"><font-awesome-icon :icon="['fas', 'toggle-off']" size="lg" /></button>
    </div>
    <div class="info">
      <span class="font-hairline">Last update: {{ asset.last_update }} </span>
      <span>{{ asset.name }}</span>
      <span>{{ asset.type }}</span>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'HiveSwitch',

    props: ['asset'],

    mounted: function() {
      this.getState();
      this.pollData();
   },

    beforeDestroy: function () {
      console.log("removing polling");
      clearInterval(this.polling);
    },

    data: function() {
      return {
        error: false,
        polling: null
      }
    },

    methods: {
      getState: function() {
        axios.get(this.asset.controller + '/state/' + this.asset.identifier)
          // force reactivity by using this.$set to add fields
          .then( response => {
            this.$set(this.asset, 'state', response.data.state);
            this.$set(this.asset, 'last_update', response.data.last_update);
          })
          .catch( error => {
            this.error = true;
            console.log(error.message);
            console.log(error.data);
          })
      },

      toggle: function() {
        var state = (this.asset.state == 'on' ? { identifier: this.asset.identifier, state: 'off', type: this.asset.type } : { identifier: this.asset.identifier, state: 'on', type: this.asset.type });
        axios.post(this.asset.controller + '/state/' + this.asset.identifier, state)
          // force reactivity by using this.$set to add fields
          .then( response => this.$set(this.asset, 'state', response.data.state))
          .catch( error => {
            this.error = true;
            console.log(error.message);
            console.log(error.data);
          })
      },

      pollData: function() {
        this.polling = setInterval( () => {
          this.error = false;
          console.log("polling");
          this.getState();
        }, 5000);
      }
    }
  }
</script>

<style scoped>

</style>
