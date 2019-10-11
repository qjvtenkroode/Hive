<template>
  <div class="flex flex-wrap w-1/4 rounded-lg m-1 p-1 border border-orange-500">
    <div class="flex w-full">
      <div class="flex w-1/3">
        <font-awesome-icon :icon="['fas', 'random']" size="lg" />
      </div>
      <div class="flex w-2/3 justify-end">
        <button v-if="asset.state === 'on'" v-on:click="toggle"><font-awesome-icon :icon="['fas', 'toggle-on']" size="lg" /></button>
        <button v-else v-on:click="toggle"><font-awesome-icon :icon="['fas', 'toggle-off']" size="lg" /></button>
      </div>
    </div>
    <div class="w-full">
      <span>Last update: timestamp </span>
    </div>
    <div class="w-full">
      <span>{{ asset.name }}</span>
    </div>
    <div class="flex w-full">
      <div class="w-3/4 rounded-full border border-green-500 text-center lowercase">
        <span>{{ asset.type }}</span>
      </div>
      <div class="flex w-1/4 justify-end">
        <span v-if="error"><font-awesome-icon :icon="['fas', 'exclamation-triangle']" /></span>
        <span v-else><font-awesome-icon :icon="['far', 'check-circle']" /></span>
      </div>
    </div>
  </div>

</template>

<script>
  export default {
    name: "HiveSwitch",

    props: ['asset'],

    mounted: function() {
      this.getState()
      
      setInterval(function() {
        this.error = false
        this.getState()
      }.bind(this), 5000);
    },

    data: function() {
      return {
        error: false
      }
    },

    methods: {
      getState: function() {
        var vm = this;
        axios.get('http://' + host + ':9000/state/' + this.asset.identifier)
          // force reactivity by using this.$set to add fields
          .then(response => this.$set(this.asset, 'state', response.data.state))
          .catch(function(error) {
            vm.error = true;
            console.log(error.message);
            console.log(error.data);
          })
      },

      toggle: function() {
        var vm = this;
        var state = (this.asset.state == 'on' ? { identifier: this.asset.identifier, state: 'off', type: this.asset.type } : { identifier: this.asset.identifier, state: 'on', type: this.asset.type });
        axios.post('http://' + host + ':9000/state/' + this.asset.identifier, state)
          // force reactivity by using this.$set to add fields
          .then(response => this.$set(this.asset, 'state', response.data.state))
          .catch(function(error) {
            vm.error = true;
            console.log(error.message);
            console.log(error.data);
          })
      }
    }
  }
</script>

<style scoped>

</style>
