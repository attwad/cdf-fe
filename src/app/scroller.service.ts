import { Injectable } from '@angular/core';

@Injectable()
export class ScrollerService {

  constructor() { }

  scrollToTop(): void {
    let toolbar = document.getElementById('toolbar');
    toolbar.scrollIntoView({block: "end", behavior: "smooth"});
  }
}
