import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LazyLoadImgComponent } from './lazy-load-img.component';

describe('LazyLoadImgComponent', () => {
  let component: LazyLoadImgComponent;
  let fixture: ComponentFixture<LazyLoadImgComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LazyLoadImgComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LazyLoadImgComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
