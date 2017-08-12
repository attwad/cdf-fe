import { CdfFePage } from './app.po';

describe('cdf-fe App', () => {
  let page: CdfFePage;

  beforeEach(() => {
    page = new CdfFePage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!');
  });
});
