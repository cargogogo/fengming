import { FengmingUiPage } from "./app.po";

describe("Fengming-ui App", () => {
  let page: FengmingUiPage;

  beforeEach(() => {
    page = new FengmingUiPage();
  });

  it("should display welcome message", () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual("Welcome to app!");
  });
});
