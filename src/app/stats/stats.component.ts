import { Component, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { StatsService, Stats } from '../stats.service';

@Component({
  selector: 'app-stats',
  templateUrl: './stats.component.html',
  styleUrls: ['./stats.component.css']
})
export class StatsComponent implements OnInit {
  stats: Stats;
  hasError = false;

  constructor(private statsService: StatsService) { }

  ngOnInit() {
    this.statsService.getStats()
    .subscribe((stats: Stats) => {
      this.stats = stats;
    },
    (err: HttpErrorResponse) => {
      console.error("error getting stats:", err);
      this.hasError = true;
    });
  }

}
