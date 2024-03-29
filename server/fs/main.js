function hideForm() {
  document.querySelectorAll('.ui').forEach(el => el.style.display = 'none')
}
function showForm(section) {
  hideForm()
  document.querySelector('#' + section + '.ui').style.display = 'block'
}
const UI = {
  decide: function(iter) {
    return new Promise((resolve, reject) => {
      iter.tasks().then(tasks => {
        document.querySelectorAll('#decide.ui .button-group button').forEach(el => el.style.display = 'none')
        showForm('decide')
        for (i = 0; i < tasks.length; i++) {
          const j = i;
          const btn = document.querySelector('#decide.ui .button-group button#btn-task-' + i)
          btn.innerText = tasks[i]
          btn.onclick = function() { resolve(j) }
          btn.style.display = 'block'
        }
      })
    })
  },
  confirm: function(prompt, task) {
    document.getElementByID('confirm-prompt').innerText = prompt
    document.getElementByID('confirm-task').innerText = task
    const btnYes = document.getElementById('btn-yes')
    const btnNo = document.getElementById('btn-no')
    showForm('confirm')
    return new Promise((resolve, reject) => {
      btnYes.onclick = function() { resolve(true) }
      btnNo.onclick = function() { resolve(false) }
    })
  },
  prompt: function(prompt) {
    document.querySelector('#prompt.ui form label').innerText = prompt
    const input = document.querySelector('#prompt.ui form input')
    input.value = ''
    showForm('prompt')
    return new Promise((resolve, reject) => {
      document.querySelector('#prompt.ui form button#prompt-enter').onclick = function() {
        resolve(input.value)
      }
    })
  },
  list: function(tasks) {
    const ul = document.querySelector('#list.ui ul')
    ul.replaceChildren()
    for (i = 0; i < tasks.length; i++) {
      const li = document.createElement('li')
      li.innerText = tasks[i]
      ul.appendChild(li)
    }
    showForm('list')
  }
}
function init () {
  document.querySelectorAll('.ui').forEach(el => el.style.display = 'none')
  const taskStream = localStorage.getItem("tasks") || "[]"
  const tasks = JSON.parse(taskStream)
  newApp(tasks, UI).then(app => {
    window.app = app
    app.peek().then(task => document.getElementById('current-task').innerText = task)
    document.getElementById('nav-complete').onclick = async function() {
      try {
        await app.complete(false)
      } catch (e) {
        console.log("couldn't complete", e)
      }
      var task
      try {
        task = await app.peek()
        document.getElementById('current-task').innerText = task
      } catch (e) {
        if (e.message === 'no tasks to do') {
          task = "All caught up!"
        } else {
          console.log("couldn't peek", e.message)
          return
        }
      }
      localStorage.setItem('tasks', JSON.stringify(app.tasks()))
      hideForm()
    }
    document.getElementById('nav-add').onclick = async function() {
      try {
        await app.add("")
      } catch (e) {
        console.log("error adding", e)
      }
      try {
        task = await app.peek()
      } catch (e) {
        console.log("error peeking", e)
        return
      }
      document.getElementById('current-task').innerText = task
      const tasks = app.tasks()
      if (typeof tasks === 'undefined') {
        return
      }
      localStorage.setItem('tasks', JSON.stringify(app.tasks()))
      hideForm()
    }
    document.getElementById('nav-resort').onclick = async function() {
      await app.resort()
      task = await app.peek()
      document.getElementById('current-task').innerText = task
      localStorage.setItem('tasks', JSON.stringify(app.tasks()))
      hideForm()
    }
    document.getElementById('nav-resort-all').onclick = async function() {
      await app.resortAll()
      task = await app.peek()
      document.getElementById('current-task').innerText = task
      localStorage.setItem('tasks', JSON.stringify(app.tasks()))
      hideForm()
    }
    document.getElementById('nav-list-all').onclick = async function() {
      await app.list()
    }
  }).catch(e => {
    console.log("tasks", tasks)
    console.log("UI", UI)
    console.log(e)
    return
  })
}
