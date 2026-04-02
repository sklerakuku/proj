<script setup>
</script>

<template>
  <header>
    <h1>Tasks Workflow Manager</h1>
    <div class="user">
      <h3 class="name">{{  }}</h3>
      <img src="" alt="user picture" class="user-picture">
    </div>
  </header>

  <main>
     <router-view></router-view>
  </main>

  <footer>
    <p> &copy; {{ new Date().getFullYear() }} =^-^=</p>
  </footer>
</template>

<style scoped>
header {
  position: sticky;
  top: 0;
  padding: 1rem 15%;
  background: var(--color-background);
  color: var(--color-text);
  border-bottom: 2px solid var(--color-text);
  transform: rotate(-0.5deg);
  font-family: var(--font-1);
  font-size: 1.5rem;
  border-radius: 255px 15px 225px 15px/15px 225px 15px 255px;
  backdrop-filter: blur(10px);
  z-index: 100;
  display: flex;
  justify-content: space-between;
}

header h1 {
  margin: 0;
  font-size: inherit;
}

main {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 140px);
  padding: 2rem;
}

footer {
  font-family: var(--font-1);
  text-align: center;
  margin: 3em;
  color: var(--color-text);
  opacity: 0.7;
  font-size: 0.75rem;
  transform: rotate(0.5deg);
}

@media (max-width: 768px) {
  header {
    padding: 1rem 5%;
    font-size: 1.2rem;
  }
  
  main {
    padding: 1rem;
  }
}
</style>