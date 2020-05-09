import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrashDetailsComponent } from './trash-details.component';

describe('TrashDetailsComponent', () => {
  let component: TrashDetailsComponent;
  let fixture: ComponentFixture<TrashDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TrashDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrashDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
