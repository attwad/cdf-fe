import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { LessonsComponent } from './lessons/lessons.component';
import { SearchComponent } from './search/search.component';
import { LessonsService } from './lessons.service';
import { ScrollerService } from './scroller.service';

import { AppRoutingModule }     from './app.routing.module';

import {MdToolbarModule} from '@angular/material';
import {MdCardModule} from '@angular/material';
import {MdButtonModule} from '@angular/material';
import {MdInputModule} from '@angular/material';
import {MdSlideToggleModule} from '@angular/material';

@NgModule({
  declarations: [
    AppComponent,
    LessonsComponent,
    SearchComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    MdToolbarModule,
    MdCardModule,
    MdButtonModule,
    MdInputModule,
    MdSlideToggleModule
  ],
  providers: [LessonsService, ScrollerService],
  bootstrap: [AppComponent]
})
export class AppModule { }
