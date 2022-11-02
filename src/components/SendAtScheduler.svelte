<script lang="ts">
  import Button from "./Button.svelte";
  import { TIME_RANGE } from '../constants'
  import type { TimeRange } from "../models";
  import { onMount } from "svelte";

  export let value: Date
  let _date: string
  let _time: string

  let dates: Date[] = []

  onMount(() => {
    function addDays(date: Date, daysToAdd: number) {
      date.setDate(date.getDate() + daysToAdd);
      return date;
    }

    let d: Date[] = []
    for(let i = 0; i < 8; i++) {
      d.push(addDays(new Date(), i))
    }
    dates = d
  })

  function setDate(date: Date) {
    value = date
    _date = date.toLocaleDateString()
    _time = date.toLocaleTimeString()
  }

  function setTimeRange(range: TimeRange) {
    let min = range.start
    let max = range.end
    let seconds = Math.floor(Math.random() * (max - min + 1) + min)
    const newDate = new Date(
      value.getFullYear(),
      value.getMonth(),
      value.getDate(),
      0,
      0,
      0,
      0
    )
    newDate.setSeconds(seconds)
    setDate(newDate)
  }

  function setSelectedDate(date: Date) {
    let newDate = new Date(
      date.getFullYear(),
      date.getMonth(),
      date.getDate(),
      value.getHours(),
      value.getMinutes(),
      0,
      0
    )
    setDate(newDate)
  }

</script>

<div class="mb-2 p-2">
  <h2>Send at: {value.toLocaleDateString()} {value.toLocaleTimeString()}</h2>
  <input type="date" bind:value={_date} />
  <input type="time" bind:value={_time}/>
  <h3>Date</h3>
  <div class="grid grid-cols-2 gap-2 mb-2">
    {#each dates as d}
    <Button title={d.toLocaleDateString()} onClick={() => setSelectedDate(d)} />
    {/each}
  </div>
  <h3>Time</h3>
  <div class="grid grid-cols-2 gap-2">
    <Button title="Night (12-6am)" onClick={() => setTimeRange(TIME_RANGE.NIGHT)} />
    <Button title="Morning" onClick={() => setTimeRange(TIME_RANGE.MORNING)} />
    <Button title="Afternoon" onClick={() => setTimeRange(TIME_RANGE.AFTERNOON)} />
    <Button title="Evening" onClick={() => setTimeRange(TIME_RANGE.EVENING)} />
  </div>
</div>


