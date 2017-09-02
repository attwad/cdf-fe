import { Component, OnInit } from '@angular/core';

import { StatsService, Stats } from '../stats.service';

@Component({
  selector: 'app-stats',
  templateUrl: './stats.component.html',
  styleUrls: ['./stats.component.css']
})
export class StatsComponent implements OnInit {
  private stats: Stats;

  constructor(private statsService: StatsService) { }

  ngOnInit() {
    this.statsService.getStats()
    .subscribe((stats: Stats) => {
      this.stats = stats;
    })
  }

}
