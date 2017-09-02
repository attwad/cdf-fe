import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { LessonsComponent } from './lessons/lessons.component';
import { SearchComponent } from './search/search.component';
import { LessonsService } from './lessons.service';
import { ScrollerService } from './scroller.service';
import { StatsService } from './stats.service';

import { AppRoutingModule } from './app.routing.module';

import {MdToolbarModule} from '@angular/material';
import {MdCardModule} from '@angular/material';
import {MdButtonModule} from '@angular/material';
import {MdInputModule} from '@angular/material';
import {MdSlideToggleModule} from '@angular/material';
import {MdProgressBarModule} from '@angular/material';
import {MdProgressSpinnerModule} from '@angular/material';
import {MdChipsModule} from '@angular/material';

import { AboutComponent } from './about/about.component';
import { StatsComponent } from './stats/stats.component';

@NgModule({
  declarations: [
    AppComponent,
    LessonsComponent,
    SearchComponent,
    AboutComponent,
    StatsComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    AppRoutingModule,
    MdToolbarModule,
    MdCardModule,
    MdButtonModule,
    MdInputModule,
    MdSlideToggleModule,
    MdProgressBarModule,
    MdProgressSpinnerModule,
    MdChipsModule
  ],
  providers: [LessonsService, ScrollerService, StatsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
