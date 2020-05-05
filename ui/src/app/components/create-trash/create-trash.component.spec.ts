import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateTrashComponent } from './create-trash.component';

describe('CreateTrashComponent', () => {
  let component: CreateTrashComponent;
  let fixture: ComponentFixture<CreateTrashComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateTrashComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateTrashComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
