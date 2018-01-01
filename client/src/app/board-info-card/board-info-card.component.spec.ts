import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BoardInfoCardComponent } from './board-info-card.component';

describe('BoardInfoCardComponent', () => {
  let component: BoardInfoCardComponent;
  let fixture: ComponentFixture<BoardInfoCardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [BoardInfoCardComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BoardInfoCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
